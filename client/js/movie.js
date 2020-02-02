const Item = ({ Title, Year, Poster}) => `
<div class="movie-card">
    <div class="card-poster">
        <img src="${Poster}" class="float-left" />
        <div class="title-year">
            <a class="title">${Title}</a>
            <a class="year">${Year}</a>
        </div>
    </div>
</div>`

function presentMoviesPopup(movies) {
    var divMovies = $('#moviesSearchResultPopup');
    divMovies.html(movies.Search.map(Item).join(''));
    closeMovieSearchForm();
    openMovieCardsDiv();
}
