var apiUrl = './api/v1/'

Vue.config.delimiters = ['${', '}'];
Vue.config.debug = true;

/**
 * Application
 * @type {Object}
 */
var App = {
	components: {
		dashboard: Vue.extend({
			template: '<p>Hello world</p>'
		}),

		page: {
			/** 
			* View all page list
			*/
			list: Vue.extend({
				template: ''
				+ '<div class="uk-grid">'
					+ '<div class="uk-width-8-10">'
						+ '<h1><i class="uk-icon-rss"></i> Page of Category / Subcategory List</h1>'
						+ '<a class="add" v-link="{path: \'/page/add\'}"><i class="uk-icon-plus"></i></a>'
					+ '</div>'
					+ '<div class="uk-width-2-10">'
						+ '<form class="uk-form">'
							+ '<input type="search" placeholder="Search" v-model="search">'
						+ '</form>'
					+ '</div>'
				+ '</div>'
				+ '<div class="uk-grid">'
					+ '<div class="uk-width-1-1">'
					+ '<table class="uk-table uk-table-hover uk-table-striped uk-table-condensed">'
						+ '<thead>'
							+ '<tr>'
								+ '<th style="width:40%">URL</th>'
								+ '<th>Description</th>'
								+ '<th style="width:50px">Action</th>'
							+ '<tr>'
						+ '</thead>'
						+ '<tbody>'
							+ '<tr v-for="item in list|filterBy search">'
								+ '<td>${item.url}</td>'
								+ '<td>${item.desc}</td>'
								+ '<td><a class="remove" @click="remove(item.id)"><i class="uk-icon-remove"></i></a> <a class="edit"><i class="uk-icon-pencil"></i></a></td>'
							+ '</tr>'
						+ '</tbody>'
					+ '</table>'
					+ '</div>'
				+ '</div>'
				,
				route: {
					canReuse: false,
					waitForData: true,
					data: function () {
						return this.fetch()
					},
				},
				methods: {
					fetch: function (fn) {
						return $.getJSON(apiUrl.concat('list/page')).then(function (response) {
							if (response.error) {
								alert(response.message)
								return
							}

							return {list: response.list}
						})
					},

					remove: function (id) {
						var self = this

						if (confirm('Are you sure?')) {
							console.log('sure')
							$.ajax({
								url: apiUrl.concat('page', '/', id),
								type: 'DELETE',
								dataType: 'JSON',
								success: function (response) {
									if (response.error) {
										alert(response.message)
										return
									}

									self.$router.go('/')
									setTimeout(function () {
										self.$router.go('/page')
									}, 10)
								}
							});
						}
					}
				}
			}),


			/**
			 * Add New Page
			 * @type {Object}
			 */
			add: Vue.extend({
				template: 'adds'
			}),


			/**
			 * Edit page
			 * @type {Object}
			 */
			edit: Vue.extend({
				template: ''
			})
		},

		item: Vue.extend({
			template: '<p>Hello world</p>'
		}),
	},

	// Main App, we need non instance
	init: Vue.extend({
		data: function () {
			return {
				navigation: [
					{icon: 'dashboard', 'path': '/'},
					{icon: 'rss', 'path': '/page'},
					{icon: 'bookmark-o', 'path': '/item'},
				]
			}
		}
	}),

	// Router
	router: new VueRouter({
		linkActiveClass: 'active'
	})
}


/**
 * Routers
 * @type {Object}
 */
App.router.map({
	'/': {
		component: App.components.dashboard
	},

	'/page': {
		component: App.components.page.list
	},

	'/page/add': {
		component: App.components.page.add
	},

	'/page/:id': {
		component: App.components.page.edit
	},

	'/item/:id': {
		component: App.components.item
	}
});

// Start
App.router.start(App.init, '#app')