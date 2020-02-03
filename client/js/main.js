function onDocumentReady() {
    $('#addMovieButton').on('click', triggerOpenDialog)
}

function openMovieSearchForm() {
    $("#movieSearchDiv").style.display = "block";
}

function closeMovieSearchForm() {
    $("#movieSearchDiv").style.display = "none";
}

function openMovieCardsDiv() {
    $("#moviesSearchResultPopup").style.display = "block";
}

function triggerOpenDialog() {
    $('#fileInput').trigger('click');
}
