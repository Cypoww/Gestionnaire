console.log("Lauching script")

var books = []

//const dataImg = "https://club-stephenking.fr/wp-content/uploads/2021/07/apres-roman-de--albinmichel-novembre-2021.jpg"

fetch('http://localhost:8080/livres')
  .then(response => response.json())
  .then(data => {
    // utilisez les données ici
    parsedData = JSON.parse(data);
    for (var i = 0; i < parsedData.length; i++) {
      active = false
      if (i == 0) {
        active = true
      } 

      addCard(parsedData[i], active)
      console.log(parsedData[i])
      books.push(parsedData[i]);
    }
  })
  .catch(error => {
    console.error(error);
});

function addCard(book, active) {
  // create a new div element
  rootDiv = document.createElement("div");
  rootDiv.classList.add("root");

  if (active == true) {
    rootDiv.classList.add("active");
  }

  cardDiv = document.createElement("div");
  cardDiv.classList.add("card");

  img = document.createElement("img");
  img.classList.add("card-img-top");
  img.src = book["img"]

  cardBody = document.createElement("div");
  cardBody.classList.add("card-body");

  h5Div = document.createElement("h5");
  h5Div.classList.add("card-title")
  h5Div.textContent = book["title"]

  pDiv = document.createElement("p");
  pDiv.classList.add("card-text")
  pDiv.textContent = book["author"]

  h6Div = document.createElement("h6");
  h6Div.classList.add("card-year")
  h6Div.textContent = book["year"]



  rootDiv.appendChild(cardDiv);
  cardDiv.appendChild(img)
  cardDiv.appendChild(cardBody)
  cardDiv.appendChild(h5Div)
  cardDiv.appendChild(pDiv)
  cardDiv.appendChild(h6Div)

  // add the newly created element and its content into the DOM
  const currentDiv = document.getElementById("entrypoint");
  currentDiv.after(rootDiv);
}
