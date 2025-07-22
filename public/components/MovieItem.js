export class MovieItemComponent extends HTMLElement {

    constructor(movie) {
        super();
        this.movie = movie;
        // console.log(this.movie.title)
    }

    connectedCallback() {
        const url = "/movies/" + this.movie.id;
        this.innerHTML = `
            <a href="${url}">
                <article>
                    <img src="${this.movie.poster_url}" 
                        alt="${this.movie.title} Poster">
                    <p>${this.movie.title} (${this.movie.release_year})</p>
                </article>
            </a>
        `;

        this.querySelector('a').addEventListener('click', e => {
            e.preventDefault();
            window.app.Router.go(url)
        })
    }


}

customElements.define('movie-item', MovieItemComponent);