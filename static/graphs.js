// Line graphs using D3

// Draw a line graph
function lineGraph(div, dates, series, labels, colours) {

    // Empty the canvas
    let cvs = d3.select(div);
    cvs.html(null);

    // Margins
    let margin = { top: 20, right: 30, bottom: 30, left: 40 },
        width = cvs.node().getBoundingClientRect().width,
        height = cvs.node().getBoundingClientRect().height;
    
    // Get maximum values from data series, depending on the current graph
    let yMax = 0;
    for ( let i = 0; i < series.length; i++ ) {
        let max = d3.max(series[i]);
        if ( max > yMax ) 
            yMax = max;
    }
    
    // Append the svg object to the body of the page
    var svg = cvs
        .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
        .append("g")
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    // Function to parse a yyyy-mm-dd date
    var pd = d3.timeParse("%Y-%m-%d");

    // Add X axis --> it is a date format
    var x = d3.scaleTime()
        .domain(d3.extent(dates, function(d) { return pd(d); }))
        .range([ 0, width ]);
    svg.append("g")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    // Add Y axis
    var y = d3.scaleLinear()
        .domain([0, yMax])
        .range([height, 0]);
    svg.append("g")
        .call(d3.axisLeft(y));

    // Add each line for the currently selected graph
    for ( let i = 0; i < series.length; i++ ) {
        svg.append("path")
                .datum(series[i])
                .attr("fill", "none")
                .attr("stroke", colours[i])
                .attr("stroke-width", 2)
                .attr("d", d3.line()
                    .x(function(d, i) { return x(pd(dates[i])) })
                    .y(function(d) { return y(d) })
                );
    }

    // Draw legend
    let xl = margin.left + width - 150,
        yl = margin.top - 20;
    for ( let i = 0; i < labels.length; ++i ) {
        svg.append("line")
                .attr("x1", xl).attr("x2", xl+20)
                .attr("y1", yl).attr("y2", yl)
                .attr("stroke", colours[i]);
        svg.append("text").text(labels[i])
                .attr("x", xl+30).attr("y", yl+5);
        yl += 15;
    }
}

// Waterfall bar graph: last number is full height, the other ones are positioned
// vertically to stack.
function waterfallGraph(div, labels, values, colours) {

     // Empty the canvas
     let cvs = d3.select(div);
     cvs.html(null);

     // Margins
     let margin = { top: 20, right: 30, bottom: 30, left: 60 },
         width = cvs.node().getBoundingClientRect().width,
         height = cvs.node().getBoundingClientRect().height;

    // Append the svg object to the body of the page
    let svg = cvs.append("svg")
                    .attr("width", width + margin.left + margin.right)
                    .attr("height", height + margin.top + margin.bottom)
                .append("g")
                    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

     // Get the height of the stacked bars, or the last bar (take the max)
     let yMax = 0;
     for ( let i = 0; i < values.length - 1; ++i ) {
        yMax += values[i];
     }
     if ( values[values.length-1] > yMax ) {
        yMax = values[values.length-1];
     }
 
     // Add Y axis
    let y = d3.scaleLinear()
        .domain([0, yMax])
        .range([0, height]);
    svg.append("g")
        .call(d3.axisLeft(y));

    // TODO: Categorical X axis, with labels
    svg.append("line")
        .attr("x1", 0)
        .attr("x2", width)
        .attr("y1", height)
        .attr("y2", height)
        .attr("stroke", "black");

    // Draw each block
    let x = margin.left, // x position of the first bar
        dx = (width - margin.left - margin.right) / values.length, // width of each bar
        bot = height; // bottom of the first bar
    for ( let i = 0; i < values.length; ++i ) {

        // Last bar starts at bottom again
        if ( i == values.length - 1 ) {
            bot = height;
        }

        // Draw block
        let h = y(values[i]),
            top = bot - h;
        svg.append("rect")
            .attr("x", x+10)
            .attr("y", top)
            .attr("width", dx-20)
            .attr("height", h)
            .style("fill", colours[i]);

        // Horizontal line connecting this block with the previous
        if ( i > 0 ) {
            svg.append("line")
            .attr("x1", x - 8)
            .attr("x2", x + 8)
            .attr("y1", bot)
            .attr("y2", bot)
            .attr("stroke", "gray");
        }

        // Draw data value at top of block
        svg.append("text").text(values[i])
            .attr("x", x + dx/2)
            .attr("y", top - 2)
            .style("font-size", 12)
            .attr("text-anchor", "middle");

        // Draw label below block, centred
        svg.append("text").text(labels[i])
            .attr("x", x + dx/2)
            .attr("y", height+15)
            .style("font-size", 12)
            .attr("text-anchor", "middle");
        
        // For a stacked waterfall, make the bottom of the next block
        // alight with the top of this one
        bot = top;

        // X position of the next block
        x += dx;
     }

}
