import { routes } from "./Routes.js";

export const Router = {
    init: () => {
        window.addEventListener('popstate', () => {
            Router.go(location.pathname, false);
        })
        // enhance current links in the document
        document.querySelectorAll('a.navlink').forEach(a => {
            a.addEventListener('click', event => {
                event.preventDefault();
                const href = a.getAttribute('href');
                Router.go(href);
            })
        })
        // go to the initial route
        Router.go(location.pathname, location.search)
    },
    go: (route, addToHistory=true) => {
        if (addToHistory) {
            history.pushState(null, "", route);
        }

        let pageElement = null;

        let routePath = route.includes("?") ? route.split("?")[0] : route;
        // console.log(routePath);
        for (const r of routes) {
            if (typeof r.path === 'string' && r.path === routePath) {
                // console.log(routePath)
                pageElement = new r.component();
                break;
            } else if (r.path instanceof RegExp) {
                const match = r.path.exec(route);
                if (match) {
                    pageElement = new r.component();
                    const params = match.slice(1);
                    pageElement.params = params;
                    break;
                }

            }
        }

        if (pageElement == null) {
            pageElement = document.createElement("h1");
            pageElement.textContent = 'Page not found';
        } 

        // inserting a new page
        const oldPage = document.querySelector('main').firstElementChild;
        if (oldPage) oldPage.style.viewTransitionName = "old";
        pageElement.style.viewTransitionName = "new";

        function updatePage() {
            document.querySelector("main").innerHTML = '';
            document.querySelector("main").appendChild(pageElement);
        }

        if (!document.startViewTransition) {
            updatePage();
        } else {
            document.startViewTransition(() => {
                updatePage();
            })
        }
        
    }
}