var apiUrl = './api/v1/'

Vue.config.delimiters = ['${', '}'];
Vue.config.debug = true;


var Wrapper = Vue.extend({
	template: ''
	+ '<div class="uk-grid">'
		+ '<div :class="{\'uk-width-10-10\': !formSearch, \'uk-width-8-10\': formSearch}">'
			+ '<h1><i class="uk-icon-${icon}"></i> ${title}</h1>'
			+ '<a class="add" v-link="{path: \'/page/add\'}" v-if="addButton"><i class="uk-icon-plus"></i></a>'
		+ '</div>'
		+ '<div class="uk-width-2-10" v-if="formSearch">'
			+ '<form class="uk-form">'
				+ '<input type="search" placeholder="${searchLabel}" v-model="searchModel" class="uk-width-1-1">'
			+ '</form>'
		+ '</div>'
	+ '</div>'
	+ '<hr />'
	+ '<div class="uk-grid">'
		+ '<div class="uk-width-1-1">'
		+ '<slot></slot>'
		+ '</div>'
	+ '</div>'
	,
	props: {
		icon: {
			required: true,
			type: String,
			default: ''
		},
		title: {
			required: true,
			type: String,
			default: 'Title'
		},
		addButton: {
			type: Boolean,
			default: false
		},
		formSearch: {
			type: Boolean,
			default: false
		},
		searchLabel: {
			type: String,
			default: 'Search'
		},
		searchModel: {
			type: String,
			default: ''
		}
	}
});


/**
 * Application
 * @type {Object}
 */
var App = {
	dashboard: Vue.extend({
		template: '<p>Hello world</p>'
	}),

	page: {
		/** 
		* View all page list
		*/
		list: Vue.extend({
			components: {wrapper: Wrapper},
			template: ''
			+ '<wrapper icon="rss" title="Watch List" :form-search="true" :search-model.sync="search" search-label="Search Watcher" :add-button="true">'
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
							+ '<td><strong>${item.title}</strong><br /><a class="url" target="_blank" href="${item.url}">${item.url} <i class="uk-icon-external-link"></i></a></td>'
							+ '<td>${item.desc}</td>'
							+ '<td><a class="remove" @click="remove(item.id)"><i class="uk-icon-remove"></i></a> <a class="edit" v-link="{path: \'/page/\' + item.id}"><i class="uk-icon-pencil"></i></a></td>'
						+ '</tr>'
					+ '</tbody>'
				+ '</table>'
			+ '</wrapper>'
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
			components: {wrapper: Wrapper},
			template: ''
			+ '<wrapper icon="rss" title="Add Watch Page">'
			+ '<form class="uk-form uk-form-horizontal">'
				+ '<div class="uk-form-row">'
					+ '<label class="uk-form-label" for="url">Page URL</label>'
					+ '<div class="uk-form-controls">'
					+ '<input type="text" id="url" placeholder="Theme Forest Category / Subcategory Page URL" class="uk-width-5-10" v-model="url">'
					+ '</div>'
				+ '</div>'
				+ '<div class="uk-form-row">'
					+ '<label class="uk-form-label" for="title">Title Alias</label>'
					+ '<div class="uk-form-controls">'
					+ '<input type="text" id="title" placeholder="Title of current URL" class="uk-width-3-10"  v-model="title">'
					+ '</div>'
				+ '</div>'
				+ '<div class="uk-form-row">'
					+ '<label class="uk-form-label" for="desc">Description</label>'
					+ '<div class="uk-form-controls">'
					+ '<textarea id="desc" cols="30" rows="5" placeholder="Page Description" class="uk-width-3-10" v-model="desc"></textarea>'
					+ '</div>'
				+ '</div>'
				+ '<div class="uk-form-row">'
				+ '<a class="uk-button uk-button-primary uk-button-large" @click="save()"><i class="uk-icon-save"></i> Save</a>'
				+ '</div>'
			+ '</form>'
			+ '</wrapper>'
			,
			data: function () {
				return {
					url: '',
					title: '',
					desc: ''
				}
			},
			methods: {
				valid: function () {
					return (/themeforest\.net/g.test(this.url) && this.url !== '') && this.title !== ''
				},

				save: function () {
					var self = this
					if (self.valid()) {
						$.post(apiUrl.concat('page'), {
							url: self.url,
							title: self.title,
							desc: self.desc
						}).then(function (response) {
							if (response.error) {
								alert(response.message)
							}

							self.$router.go('/page')
						})
					} else {
						alert("Go home you're drunk, is URL / Title valid?")
					}
				}
			}
		}),


		/**
		 * Edit page
		 * @type {Object}
		 */
		edit: Vue.extend({
			components: {wrapper: Wrapper},
			template: ''
			+ '<wrapper icon="rss" :title="\'Edit Watch Page &quot;\'+ title +\'&quot;\'">'
			+ '<form class="uk-form uk-form-horizontal">'
				+ '<div class="uk-form-row">'
					+ '<label class="uk-form-label" for="url">Page URL</label>'
					+ '<div class="uk-form-controls">'
					+ '${url}'
					+ '</div>'
				+ '</div>'
				+ '<div class="uk-form-row">'
					+ '<label class="uk-form-label" for="title">Title Alias</label>'
					+ '<div class="uk-form-controls">'
					+ '<input type="text" id="title" placeholder="Title of current URL" class="uk-width-3-10"  v-model="title">'
					+ '</div>'
				+ '</div>'
				+ '<div class="uk-form-row">'
					+ '<label class="uk-form-label" for="desc">Description</label>'
					+ '<div class="uk-form-controls">'
					+ '<textarea id="desc" cols="30" rows="5" placeholder="Page Description" class="uk-width-3-10" v-model="desc"></textarea>'
					+ '</div>'
				+ '</div>'
				+ '<div class="uk-form-row">'
				+ '<a class="uk-button uk-button-primary uk-button-large" @click="save()"><i class="uk-icon-save"></i> Save Changes</a>'
				+ '</div>'
			+ '</form>'
			+ '</wrapper>'
			,
			data: function () {
				return {
					url: '',
					title: '',
					desc: ''
				}
			},
			route: {
				canReuse: false,
				waitForData: true,
				data: function () {
					return this.fetch()
				},
			},
			methods: {
				fetch: function (fn) {
					return $.getJSON(apiUrl.concat('page/', this.$route.params.id)).then(function (response) {
						if (response.error) {
							alert(response.message)
							return
						}

						return response.data
					})
				},

				save: function () {
					var self = this
					$.ajax({
						url: apiUrl.concat('page', '/', self.$route.params.id),
						type: 'PUT',
						data: {title: self.title, desc: self.desc},
						dataType: 'JSON',
						success: function (response) {
							if (response.error) {
								alert(response.message)
								return
							}

							self.$router.go('/page')
						}
					});
				}
			}
		})
	},

	item: {
		view: Vue.extend({
			template: '<p>Hello world</p>'
		})
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
		component: App.dashboard
	},

	'/page': {
		component: App.page.list
	},

	'/page/add': {
		component: App.page.add
	},

	'/page/:id': {
		component: App.page.edit
	},

	'/item/:id': {
		component: App.item.view
	}
});

// Start
App.router.start(App.init, '#app')