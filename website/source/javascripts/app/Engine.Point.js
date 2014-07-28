(function(
	Engine,
	Vector
){ 'use strict';

Engine.Point = function(id, x, y, shapeSize){
	this.id = id;

	this.shapeSize = shapeSize;
	this.ref = new Vector(x, y);

	this.pos = new Vector(
		x * shapeSize.x,
		y * shapeSize.y
	);

	this.target = this.pos.clone();
	this.pos.x  = shapeSize.x / 2;
	this.pos.y  = shapeSize.y / 2;
	this.accel  = Vector.coerce(this.accel);
	this.vel    = Vector.coerce(this.vel);

	this.stiffness = Engine.getRandomFloat(3, 6);
	this.friction  = Engine.getRandomFloat(0.15, 0.3);
};

Engine.Point.prototype = {

	radius: 1,

	stiffness : 0.5,
	// friction  : 0.00001,
	friction  : 0.01,
	threshold : 0.03,

	pos: {
		x: 0,
		y: 0
	},

	accel: {
		x: 0,
		y: 0
	},

	vel : {
		x: 0,
		y: 0
	},

	target: {
		x: 0,
		y: 0
	},

	resize: function(){
		this.target.x = this.pos.x = this.ref.x * this.shapeSize.x;
		this.target.y = this.pos.y = this.ref.y * this.shapeSize.y;
	},

	updateBreathingPhysics: function(){
		this.stiffness = 0.1;
		this.friction  = 0.05;
	},

	updateTarget: function(newSize){
		var diff;

		this.target.x = this.ref.x * newSize.x;
		this.target.y = this.ref.y * newSize.y;

		diff = Vector.sub(newSize, this.shapeSize).div(2);

		this.target.sub(diff);

		this.target.add({
			x: Engine.getRandomFloat(-8, 8),
			y: Engine.getRandomFloat(-8, 8)
		});
	},

	update: function(engine){
		var newAccel;

		newAccel = Vector.sub(this.target, this.pos)
			.mult(this.stiffness)
			.sub(Vector.mult(this.vel, this.friction));

		this.accel.set(newAccel);

		this.vel.add(this.accel);

		this.pos.add(
			Vector.mult(this.vel, engine.tick)
		);

		newAccel = null;

		return this;
	},

	draw: function(ctx, scale){
		ctx.beginPath();
		ctx.arc(
			this.pos.x  * scale,
			this.pos.y  * scale,
			this.radius * scale,
			0,
			Math.PI * 2,
			false
		);
		ctx.fillStyle = '#ffffff';
		ctx.fill();
		return this;
	}

};

})(window.Engine, window.Vector);
