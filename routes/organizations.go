package routes

import (
	"net/http"

	"github.com/Corwind/goauth/organizations"

	"github.com/gin-gonic/gin"
)

type OrganizationPostSchema struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func (env *Env) V1GetOrganizations(c *gin.Context) {
	ret, err := organizations.FetchOrganizations(env.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, &ret)
}

func (env *Env) V1PostOrganizations(c *gin.Context) {
	var form OrganizationPostSchema
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	orga := *organizations.NewOrganization(form.Name)
	orga_, err := organizations.SaveOrganization(env.DB, orga)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &orga_)
}

func (env *Env) V1GetOrganization(c *gin.Context) {
	organization_id := c.Param("id")
	orga, err := organizations.FetchOrganization(env.DB, organization_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &orga)
}
