package routers

import (
	"event_explorer/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/events", &controllers.EventController{}, "get:List")
	beego.Router("/events/:eventId", &controllers.EventController{}, "get:Details")
	beego.Router("/redirect/:eventId", &controllers.EventController{}, "get:Redirect")
	beego.Router("/api/locations/autocomplete", &controllers.APIController{}, "get:Autocomplete")
	beego.Router("/api/locations/:placeId", &controllers.APIController{}, "get:Location")
	beego.Router("/api/cache/invalidate", &controllers.APIController{}, "get:InvalidateCache")
	beego.Router("/api/cache/invalidate-all", &controllers.APIController{}, "get:InvalidateAllCache")
}
