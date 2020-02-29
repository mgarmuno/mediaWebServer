const MovieInfo = ({Filename, Episode, Session}) => `
<form class="movie-info-form">
    <input for="Filename" value="${Filename}" />
    <input for="Episode" value="${Episode}" />
    <input for="Session" value="${Session}" />
</form>`
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
const MovieSearchResultPopupID = '#movieSearchReslutPopup';
const MovieOptionsPopupID = '#movieOptionsPopup';

function showMoviesPopup(data) {
    var movies = JSON.parse(data);
    if (movies.Response === true) {
        var divMovies = $(MovieSearchResultPopupID);
        divMovies.html(movies.Search.map(Item).join(''));
        openPopup(MovieSearchResultPopupID);
    } else {
        console.log(movie.Error);
    }
}

function showOptions(responseJSONString) {
    var response = JSON.parse(responseJSONString);
    if (response.Options && response.Options.Search.length > 0) {
        var divMoviesOptions = $(MovieOptionsPopupID);
        divMoviesOptions.html(response);
        divMoviesOptions.append(response.Options.map(Item).join(''));
        openpopup(MovieOptionsPopupID);
    }
}
