(function(
	Engine,
	Point,
	Polygon,
	Vector
){

Engine.Shape = function(x, y, width, height, points, polygons){
	var i, ref, point, poly;

	this.pos = new Vector(x, y);
	this.size = new Vector(width, height);

	ref = {};
	this.points = [];
	this.polygons = [];

	for (i = 0; i < points.length; i++) {
		point = new Point(
			points[i].id,
			points[i].x * this.size.x,
			points[i].y * this.size.y,
			this.size.x,
			this.size.y
		);
		ref[point.id] = point;
		this.points.push(point);
	}

	for (i = 0; i < polygons.length; i++) {
		poly = polygons[i];
		this.polygons.push(new Polygon(
			ref[poly.points[0]],
			ref[poly.points[1]],
			ref[poly.points[2]],
			poly.color
		));
	}
};

Engine.Shape.prototype = {

	update: function(engine){
		var p;

		for (p = 0; p < this.points.length; p++)  {
			this.points[p].update(engine);
			// this.points[p].draw(this.context, scale);
		}

		for (p = 0; p < this.polygons.length; p++) {
			this.polygons[p].update(engine);
			// this.polygons[p].draw(this.context, scale);
		}
	},

	draw: function(ctx, scale){
		var p;

		ctx.save();
		ctx.translate(this.pos.x * scale, this.pos.y * scale);
		for (p = 0; p < this.polygons.length; p++) {
			this.polygons[p].draw(ctx, scale);
		}
		ctx.restore();
	}

};

})(
	window.Engine,
	window.Engine.Point,
	window.Engine.Polygon,
	window.Vector
);
