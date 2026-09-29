package routers

import (
	"event_explorer/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/events", &controllers.EventController{})
	beego.Router("/events/:eventId", &controllers.EventController{})
	beego.Router("/redirect/:eventId", &controllers.EventController{})
	beego.Router("/api/locations/autocomplete", &controllers.APIController{}, "get:Autocomplete",)
	beego.Router("/api/locations/:placeId",	&controllers.APIController{}, "get:Location",)
}