$(function(){
    var libLoading = false;
    var selectedBookID = -1;
        
    $('select').material_select();

    $('#library-field').change(function(){
        getBooks();
    });

    //Редактируем книгу
    $('.book-edit-trigger').click(function(){
        $('#book-modal').openModal({
            ready: function(){
                getBook(selectedBookID, function(data){
                    $('#book_id_field').val(data.id);
                    $('#book_name_field').val(data.name);
                    $('#author_field').val(data.authors);
                });

                loadLibraries(function(data){
                    $('select.library-field').each(function(){
                        populateLookup($(this), data);
                    });
                });

                Materialize.updateTextFields();
            }
        });
    });

    //Редактируем библиотеку
    $('.lib-edit-trigger').click(function(){
        if (!$('#library-field').val()){
            Materialize.toast('Необходимо выбрать библиотеку!', 2000);
            return;
        }

        $('#lib-modal').openModal({
            ready: function(){
                if ($('#library-field').val()){
                    
                    var lib_id = $('#library-field').val();
                    $('#lib_id_field').val(lib_id);
                    
                    var lib_name = $('#library-field>option[value="'+lib_id+'"]').text();
                    $('#lib_name_field').val(lib_name);
                }

                Materialize.updateTextFields();
            }
        });
    });
    
    //Создаем библиотеку
    $('.lib-create-trigger').leanModal();
    
    $('#lib_edit_action').click(function(){
            var lib_id = $('#lib_id_field').val();
            var lib_name = $('#lib_name_field').val();
            if (!lib_id){
                createLib(lib_name, function(data){
                    loadLibraries(fillLibraries);
                });
            }else{
                editLib(lib_id, lib_name, function(data){
                    loadLibraries(fillLibraries);
                });
            }
    });

    //Удаляем библиотеку
    $('.lib-delete-trigger').click(function(){
        var lib_id = $('#library-field').val();
        if (!lib_id){
            Materialize.toast('Необходимо выбрать библиотеку!', 2000);
            return;
        }

       deleteLib(lib_id, function(data){
          console.log(data);
          loadLibraries(fillLibraries);
       });
    });

    $('td>ul>.tool-item').each(function(){
        $(this).children('a').on('click', function(){
            var row = $(this).parents('tr').get(0);
            selectedBookID = row.cells[0].innerHTML;
        });
    });

    loadLibraries(fillLibraries);
    getBooks();

    $('#book-create-trigger').leanModal({
        ready: function(){
            $('#book_id_field').val(null);
            $('#book_name_field').val(null);
            $('#book_author_field').val(null);

            var lib_id = $('#library-field').val();
            if (lib_id){
                $('#book_library_field').val(lib_id);
                $('#book_library_field').material_select();
            }
        }
    });

    $('#book_edit_action').click(function(){
         var lib_id = $('#book_library_field').val();
         var book_id = $('#book_id_field').val();
         var name = $('#book_name_field').val();
         var authors = $('#book_author_field').val();

         if (!book_id){
             createBook(name, lib_id, authors, getBooks);
         }else{
             editBook(book_id, name, lib_id, authors, getBooks);
         }
    });
});



function populateLookup(dropdown, lookup) {
   var selectedValue = dropdown.val();
   if (!selectedValue && lookup.length > 0){
       selectedValue = lookup[0].ID;
   }

   dropdown.children().remove();
   
   lookup.forEach(function(item){
        newOpt = $("<option></option>").attr('value', item.ID).text(item.Name);
        if (!item.ID){
            newOpt.attr('disabled', 'disabled');
        }

        dropdown.append(newOpt);
   });
   
   dropdown.val(selectedValue);
   dropdown.material_select();
}

// Library requests
    function loadLibraries(func){
        $.ajax({
            url:'libs',
            dataType:'json',
            success: function(data){
                    successHandler(data, func);
                },
            error: showError
        });
    }

    function getLib(id, func){
        $.ajax({
            url:'lib/'+id,
            dataType:'json',
            success: function(data){
                    successHandler(data, func);
                },
            error: showError
        });
    }

    function createLib(name, func){
        var postLib = new FormData();
        postLib.append('name', name);
        postBook.append('_csrf', $('input[name="_csrf"]').val());

        $.ajax({
            url:'lib/',
            dataType:'json',
            cache: false,
            contentType: false,
            processData: false,
            data: postLib,
            type: 'post',
            success: function(data){
                    successHandler(data, func);
                },
            error: showError
        });
    }

    function editLib(id, name, func){
        var postLib = new FormData();
        postLib.append('name', name);
        postLib.append('id', id);
        postLib.append('_csrf', $('input[name="_csrf"]').val());

        $.ajax({
            url:'lib/',
            dataType:'json',
            cache: false,
            contentType: false,
            processData: false,
            data: postLib,
            type: 'put',
            success: function(data){
                    successHandler(data, func);
                },
            error: showError
        });
    }

    function deleteLib(id, func){
        var postLib = new FormData();
        postLib.append('id', id);
        postBook.append('_csrf', $('input[name="_csrf"]').val());

        $.ajax({
            url:'lib/',
            dataType:'json',
            cache: false,
            contentType: false,
            processData: false,
            data: postLib,
            method:'delete',
            success: function(data){
                    successHandler(data, func);
                },
            error: showError
        });
    }

// Book requests
    function createBook(name, lib, author, func){
        var postBook = new FormData();
        postBook.append("name", name);
        postBook.append("libraries", lib);
        postBook.append("authors", author);
        postBook.append('_csrf', $('input[name="_csrf"]').val());

        $.ajax({
            url:'book/',
            dataType:'json',
            cache: false,
            contentType: false,
            processData: false,
            data: postBook,
            method:'post',
            success: function(data){
                    successHandler(data, func);
                },
            error: showError
        });
    }

    function editBook(id, name, lib, author, func){
        var postBook = new FormData();
        postBook.append('name', name);
        postBook.append('id', id);
        postBook.append('libraries', lib);
        postBook.append('authors', author);
        postBook.append('_csrf', $('input[name="_csrf"]').val());


        $.ajax({
            url:'book/',
            dataType:'json',
            cache: false,
            contentType: false,
            processData: false,
            data: postBook,
            type: 'put',
            success: function(data){
                        successHandler(data, func);
                    },
            error: showError
        });
    }

    function deleteBook(id, func){
        var postBook = new FormData();
        postBook.append('id', id);
        postBook.append('_csrf', $('input[name="_csrf"]').val());

        $.ajax({
            url:'book/',
            dataType:'json',
            cache: false,
            contentType: false,
            processData: false,
            data: postBook,
            method:'delete',
            success: function(data){
                    successHandler(data, func);
                },
            error: showError
        });
    }

    function getBooks(){
        var lib_filter = $('#library-field').val();
        if (!lib_filter){
            lib_filter = -1;
        }

        $.ajax({
            url:'book/',
            type:'get',
            data:{lib: lib_filter},
            success: function(data){
                successHandler(data, refreshBooks);
            },
            error: showError
        });
    }

    function getBook(id, func){
        $.ajax({
            url:'book/'+id,
            dataType:'json',
            success: function(data){
                    successHandler(data, func);
                },
            error: showError
        });
    }

function refreshBooks(data){
    $('tbody').children().remove();
    data.books.forEach(function(book){
        var newRow = $('<tr></tr>');
        newRow.append($('<td>'+book.ID+'</td>'));
        newRow.append($('<td>'+book.Name+'</td>'));
        newRow.append($('<td>'+book.Authors+'</td>'));
        newRow.append('<td><ul class="toolbar"><li class="tool-item"><a class="btn-floating green book-edit-trigger" href="#book-modal"><i class="material-icons">mode_edit</i></a></li><li class="tool-item"><a class="btn-floating blue book-delete-trigger"><i class="material-icons">delete</i></a></li></ul></td>');
        $('tbody').append(newRow);
    });

    $('.book-edit-trigger').click(function(){
        $('#book-modal').openModal({
            ready: function(){
                getBook(selectedBookID, function(data){
                    $('#book_id_field').val(data.book.ID);
                    $('#book_name_field').val(data.book.Name);
                    $('#book_author_field').val(data.book.Authors);
                    $('#book_library_field').val(data.libraries);
                });

                loadLibraries(function(data){
                    $('select.library-field').each(function(){
                        populateLookup($(this), data.libraries);
                    });
                });

                Materialize.updateTextFields();
            }
        });
    });

    $('td>ul>.tool-item').each(function(){
        $(this).children('a').on('click', function(){
            var row = $(this).parents('tr').get(0);
            selectedBookID = row.cells[0].innerHTML;
        });
    });

    $('.book-delete-trigger').click(function(){
        deleteBook(selectedBookID, getBooks);
    });
}

function fillLibraries(data){
    var libraries = data.libraries;
    libraries.unshift({"ID": "", "Name": "Библиотека"});

    if (!libraries || libraries.length ==0){
        return;
    }
        
    $('select.library-field').each(function(){
        populateLookup($(this), libraries);
    });
}

function showError(data){
    Materialize.toast(data.responseJSON.status, 4000);
    $('input[name="_csrf"]').val(data.responseJSON._csrf);
}

function successHandler(data, func){
    $('input[name="_csrf"]').val(data._csrf);
    func(data);
}