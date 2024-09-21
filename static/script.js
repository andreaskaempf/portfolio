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
