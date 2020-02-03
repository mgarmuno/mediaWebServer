function onDocumentReady() {
    $('#addMovieButton').on('click', triggerOpenDialog)
}

function openMovieSearchForm() {
    $("#movieSearchDiv").first().show();
}

function closeMovieSearchForm() {
    $("#movieSearchDiv").first().hide();
}

function openMovieCardsDiv() {
    $("#moviesSearchResultPopup").first().show();
}

function triggerOpenDialog() {
    $('#fileInput').trigger('click');
}
