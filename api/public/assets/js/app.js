var apiUrl = './api/v1/'

Vue.config.delimiters = ['${', '}'];
Vue.config.debug = true;


var Wrapper = Vue.extend({
	template: '#wrapper',
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
			template: '#watcher-list',
			components: {wrapper: Wrapper},
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
			template: '#watcher-add',
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
			template: '#watcher-edit',
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
		/**
		 * Item List
		 * @type {Object}
		 */
		list: Vue.extend({
			template: '#item-list',
			components: {wrapper: Wrapper},
			route: {
				canReuse: false,
				data: function () {
					return this.fetch()
				},
			},
			
			data: function () {
				return {
					pagination: {
						current: 0,
						limit: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
						total: 1
					}
				}
			},

			methods: {
				fetch: function (offset) {
					var self = this
					offset = (offset)? offset: 0

					return $.getJSON(apiUrl.concat('list/item'), {offset: offset}).then(function (response) {
						if (response.error) {
							alert(response.message)
							return
						}

						self.pagination.current = offset
						self.pagination.total = Math.ceil(response.total / 100)

						var limit = []
						if (offset < self.pagination.total - 10) {
							for (var i = 0;i < 10;i++) {
								limit.push(i + offset)
							}
						} else {
							for (var i = 0;i < 10;i++) {
								limit.push(i + (self.pagination.total - 10))
							}
						}
						
						self.pagination.limit = limit

						return response
					})
				},

				getPage (offset) {
					var self = this
					self.fetch(offset).then(function (response) {
						self.list = response.list
					})
				},

				prevPage (event) {
					var self = this, offset = self.pagination.current - 1
					offset = (offset <= 0)? 0: offset

					self.fetch(offset).then(function (response) {
						self.list = response.list
					})
				},

				nextPage (event) {
					var self = this, offset = self.pagination.current + 1
					offset = (offset >= self.pagination.total)? self.pagination.total - 1: offset

					self.fetch(offset).then(function (response) {
						self.list = response.list
					})
				}
			}
		}),


		/**
		 * Item View
		 * @type {Object}
		 */
		view: Vue.extend({
			template: '#item-view',
			components: {wrapper: Wrapper},
			route: {
				canReuse: false,
				waitForData: true,
				data: function () {
					return this.fetch()
				},
			},
			methods: {
				fetch: function (fn) {
					return $.getJSON(apiUrl.concat('item/', this.$route.params.id)).then(function (response) {
						if (response.error) {
							alert(response.message)
							return
						}

						response.data.tags = _.unique(response.data.tags)
						return response.data
					})
				}
			},

			ready: function () {
				var self = this, chartData,
				chartOptions = {
					hAxis: {title: 'Time'},
					vAxis: {title: 'Sales'},
					backgroundColor: '#ffffff'
				}


				// Load google.visualization.DataTable()
				async.waterfall([
					function (next) {
						if (google.visualization) {
							chartData = new google.visualization.DataTable()
							next ()
						} else {
							google.charts.load('current', {packages: ['corechart', 'line']})
							google.charts.setOnLoadCallback(function () {
								chartData = new google.visualization.DataTable()
								next()
							})
						}
					}
				], function () {
					chartData.addColumn('string', 'X')
					chartData.addColumn('number', 'Sales')

					var sales = [], totalSales = _.last(self.sales).value
					_.each(self.sales, function (item, index) {
						sales.push([moment.unix(item.date).format("DD MM YYYY"), totalSales - item.value])
					})
					chartData.addRows(sales)

					var chart = new google.visualization.LineChart(document.getElementById('chart'))
					chart.draw(chartData, chartOptions);
				})
			}
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
					{icon: 'search', 'path': '/search'},
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

	'/item': {
		component: App.item.list
	},

	'/item/:id': {
		component: App.item.view
	}
});

// Start
App.router.start(App.init, '#app')