// JavaScript for the Portfolio Management System

// Handle events to change tabs on the Stock page
// Source: https://codepen.io/t7team/pen/ZowdRN
function openTab(evt, tabName) {
	
  	var i, x, tablinks;

	// Set all content areas invisible
  	var divs = document.getElementsByClassName("content-tab");
  	for ( var i = 0; i < divs.length; i++ ) {
      	divs[i].style.display = "none";
  	}

	// Make all tabs inactive
	// Warning: assumes that is-active is preceded by space, i.e.,
	// not the first or only class
  	tabs = document.getElementsByClassName("tab");
  	for ( var i = 0; i < tabs.length; i++ ) {
      	tabs[i].className = tabs[i].className.replace(" is-active", "");	
	}
	
	// Make the selected tab active, and the selected content visible
  	document.getElementById(tabName).style.display = "block";
  	evt.currentTarget.className += " is-active";
}

// Fetch prices and draw graph
async function get_prices(sid) {
	
	try {
    		// Get data
      	const response = await fetch("/get_prices/" + sid);
      	if (!response.ok) {
        		throw new Error(`Response status: ${response.status}`);
      	}

      	// Parse JSON
      	const data = await response.json();
		
      
      	// Convert to lists of dates and values
      	var dates = [], prices = [];
      	for ( var i = 0; i < data.length; ++i ) {
        		let d = data[i].Date;
        		if ( d.length > 10 )   // Convert "2018-10-01T00:00:00Z" to just date
          		d = d.substr(0, 10);
        		dates.push(d);
        		prices.push(data[i].Price);
      	}
      
      	// Show graph
      	lineGraph("#graph", dates, [prices], ["Price"], ["red"]);
      
    } catch (error) {
      	console.error(error.message);
    }
}