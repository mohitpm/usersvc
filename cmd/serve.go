package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	cfg "github.com/mohitpm/usersvc/config"
	"github.com/mohitpm/usersvc/db"
	httpHandler "github.com/mohitpm/usersvc/handler/delivery/http"
	"github.com/mohitpm/usersvc/pkg/saml"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the user service",
	Long:  `Start the user service`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

}

func startServer() {
	config := cfg.LoadConfig()

	err := db.RunMigrations(&config.Database)
	if err != nil {
		panic(err)
	}

	samlSP, err := saml.NewSP(config.SAML)
	if err != nil {
		panic(err)
	}

	idpMetadataURL, err := url.Parse(config.SAML.IDP.MetadataURL)
	if err != nil {
		panic(err)
	}

	// register with the service provider
	spMetadataBuf, _ := xml.MarshalIndent(samlSP.ServiceProvider.Metadata(), "", "  ")
	spURL := *idpMetadataURL
	spURL.Path = "/services/sp"
	resp, err := http.Post(spURL.String(), "text/xml", bytes.NewReader(spMetadataBuf))
	if err != nil {
		panic(err)
	}
	if err := resp.Body.Close(); err != nil {
		panic(err)
	}

	authHandler := httpHandler.NewAuthHandler()

	e := echo.New()
	e.Any("/saml/*", func(ctx echo.Context) error {
		samlSP.ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	})

	authGroup := e.Group("")
	authGroup.Use(echo.WrapMiddleware(samlSP.RequireAccount))
	authGroup.GET("/sign-in", authHandler.SignIn)

	address := fmt.Sprintf(":%d", config.App.Port)
	e.Logger.Fatal(e.Start(address))
}
