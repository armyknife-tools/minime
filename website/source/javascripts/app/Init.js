(function(
	Engine
){

var isIE = (function(){

	var undef,
	v = 3,
	div = document.createElement('div'),
	all = div.getElementsByTagName('i');

	while (
		div.innerHTML = '<!--[if gt IE ' + (++v) + ']><i></i><![endif]-->',
			all[0]
	);

	return v > 4 ? v : undef;

}());

// isIE = true;

var Init = {

	start: function(){
		var id = document.body.id.toLowerCase();

		if (this.Pages[id]) {
			this.Pages[id]();
		}
	},

	generateAnimatedLogo: function(){
		var container, x, block;

		container = document.createElement('div');
		container.className = 'animated-logo';

		for (x = 1; x < 5; x++) {
			block = document.createElement('div');
			block.className = 'white-block block-' + x;
			container.appendChild(block);
		}

		return container;
	},

	initializeEngine: function(){
		var jumbotron = document.getElementById('jumbotron'),
			content   = document.getElementById('jumbotron-content'),
			tagLine   = document.getElementById('tag-line'),
			canvas, galaxy;

		if (!jumbotron) {
			return;
		}

		galaxy = document.createElement('div');
		galaxy.id = 'galaxy-bg';
		galaxy.className = 'galaxy-bg';
		jumbotron.appendChild(galaxy);

		content.appendChild(
			Init.generateAnimatedLogo()
		);

		canvas = document.createElement('canvas');
		canvas.className = 'terraform-canvas';

		jumbotron.appendChild(canvas);
		window.engine = new Engine(canvas, galaxy, tagLine);
	},

	Pages: {
		'page-home': function(){
			var jumbotron;
			if (isIE) {
				jumbotron = document.getElementById('jumbotron');
				jumbotron.className += ' static';
				return;
			}

			Init.initializeEngine();
		}
	}

};

Init.start();

})(window.Engine);
