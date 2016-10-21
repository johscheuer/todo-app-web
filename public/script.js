$(document).ready(function() {
  var entryContentElement = $("#todo-input");

  var appendTodoList = function(data) {
    if (data == null) {
      return
    }
    $("#Todos > tbody").empty();
    $.each(data, function(key, val) {
      $("#Todos > tbody").append('<tr><td class="col-xs-10 col-sm-10 col-md-10">'+val+'</td><td align="center" class="col-xs-2 col-sm-2 col-md-2"><input type="checkbox" name="deleteCheck" value="1"/></td></tr>');
    });
  }

  var handleSubmission = function(e) {
    e.preventDefault();
    var entryValue = entryContentElement.val()
    if (!entryValue || entryValue.length <= 0){
        entryContentElement.parent().addClass("has-error").removeClass("has-success");
        return false;
    }


    entryContentElement.val("")
    entryContentElement.parent().removeClass("has-error").addClass("has-success");
    $.getJSON("insert/todo/" + entryValue, appendTodoList);
  }

  var handleDeletion = function(e){
    e.preventDefault();

    var checkboxes = document.getElementsByName("deleteCheck");
    for (var i=0; i < checkboxes.length; i++) {
     if (!checkboxes[i].checked){
       continue
     }
     var checkbox = checkboxes[i];
     $.getJSON("delete/todo/" + $(checkbox).closest('tr').text(), appendTodoList);
    }
  }

  $("#todo-submit").click(handleSubmission);
  $("#todo-delete").click(handleDeletion);

  // Poll every second.
  (function fetchTodos() {
    $.getJSON("read/todo").done(appendTodoList).always(
      function() {
        setTimeout(fetchTodos, 10000);
      });
  })();
});

$(document).ready(function() {
   $.getJSON("version", function(data) {
    if (data == null) {
      return
    }
    $("#footer-version").text("Version: " + data["version"]);
  });
});
