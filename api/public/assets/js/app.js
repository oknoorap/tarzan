var apiUrl = './api/v1/'

Vue.config.delimiters = ['${', '}'];
Vue.config.debug = true;
Vue.transition('show', {
	enterClass: 'slideInDown',
	leaveClass: 'slideOutUp'
})

function renderDate (date) {
	if (!date) return
	date = date.substr(0, 10).split('-')
	return date[2].concat('/', date[1], '/', date[0])
}

function orderDate (data) {
	var timeline = {};
	Object.keys(data).sort(function(a, b){
		return moment(a, 'DD/MM/YYYY').toDate() - moment(b, 'DD/MM/YYYY').toDate()
	}).forEach(function(key) {
		timeline[key] = data[key]
	});

	return timeline
}

/**
 * Application
 * @type {Object}
 */
var App = {}
App.loader = Vue.extend({
	template: '#loader'
})

/**
 * Application dashboard
 * @type {void}
 */
App.dashboard = Vue.extend({
	template: '#dashboard',
	route: {
		canReuse: false,
		waitForData: false
	},
	components: {loader: App.loader},
	data: function () {
		return {
			tab: {
				'market-value': 'today',
				'categories-value': 'today',
				'group-value': 'today',
				'bestselling': 'today'
			},

			selectedCategory: '',
			selectedGroup: '',
			categories: [],
			groups: [],
			bestselling: null,
			loader: {
				bestselling: false,
				'tag-stats': false,
				'market-value': false,
				'categories-value': false,
				'group-value': false
			},
			canvas: {
				'tag-stats': {},
				'market-value': {},
				'categories-value': {},
				'group-value': {}
			}
		}
	},
	methods: {
		/**
		 * Get Categories List
		 * @return {Object}
		 */
		getCategories: function () {
			var self = this
			$.getJSON(apiUrl.concat('list/category'))
			.then(function (response) {
				if (response.error) {
					alert(response.message)
					return
				}

				self.categories = response.list
				self.selectedCategory = self.categories[0]
			})
		},

		/**
		 * Get Group List
		 * @return {Object}
		 */
		getGroups: function () {
			var self = this
			$.getJSON(apiUrl.concat('list/subscribe/group'))
			.then(function (response) {
				if (response.error) {
					alert(response.message)
					return
				}

				self.groups = response.list
				if (self.groups.length>0) {
					self.selectedGroup = self.groups[0].id
				}
			})
		},

		/**
		 * Get Market Value (All categories) by date
		 * @param  {string} date
		 * @return {void}
		 */
		getMarketValue: function (chart, date) {
			var self = this, canvas = self.canvas[chart]

			canvas.context = document.getElementById(chart).getContext('2d')
			canvas.data = {
				labels: [],
				datasets: [{
					label: "Revenue",
					fill: true,
					lineTension: 0.1,
					backgroundColor: "rgba(26, 188, 156, 0.4)",
					borderColor: "rgba(26, 188, 156,1.0)",
					borderCapStyle: 'butt',
					borderDash: [],
					borderDashOffset: 0.0,
					borderJoinStyle: 'miter',
					pointBorderColor: "rgba(22, 160, 133,1.0)",
					pointBackgroundColor: "rgba(26, 188, 156,1.0)",
					pointBorderWidth: 2,
					pointHoverRadius: 5,
					pointHoverBackgroundColor: "#ffffff",
					pointHoverBorderColor: "rgba(22, 160, 133,1.0)",
					pointHoverBorderWidth: 2,
					pointRadius: 5,
					pointHitRadius: 10,
					data: []
				}]
			}

			var endpoint = apiUrl.concat('dashboard/stats/market?date=', date)
			if (chart === 'categories-value') {
				if (self.selectedCategory !== '') {
					endpoint = endpoint.concat('&category=', self.selectedCategory)
				}
			} else if (chart === 'group-value') {
				if (self.selectedGroup !== '') {
					endpoint = endpoint.concat('&group=', self.selectedGroup)
				}
			}

			$.getJSON(endpoint).then(function (response) {
				if (response.error) { alert(response.message); return }

				// Push actual data to chart
				_.each(orderDate(response.data), function (item, date) {
					canvas.data.labels.push(date)
					canvas.data.datasets[0].data.push(item.sales * item.price)
				})

				// Clear and Render Canvas
				self.loader[chart] = false
				if (canvas.chart) canvas.chart.destroy()
				canvas.chart = new Chart(canvas.context, {
					type: 'line',
					data: canvas.data,
					options: {
						 tooltips: {
						 	callbacks: {
						 		label: function(item, data) {
						 			return '$' + item.yLabel.toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,')
						 		}
						 	}
						 }
					}
				})
			})
		},


		/**
		 * View market value
		 * @param  {String} date
		 * @return {void}
		 */
		viewMarketValue: function (chart, date) {
			this.getMarketValue(chart, date)
			this.loader[chart] = true
			this.tab[chart] = date
		},

		/**
		 * Tags Stats
		 * @type {Object}
		 */
		getTagStats: function () {
			var self = this, canvas = self.canvas['tag-stats']

			canvas.context = document.getElementById('tag-stats').getContext('2d')
			canvas.data = {
				labels: [],
				borderWidth: 1,
				datasets: [{
					data: [],
					backgroundColor: ['#2ecc71', '#3498db', '#9b59b6', '#f1c40f', '#e67e22', '#e74c3c', '#ecf0f1', '#95a5a6', '#1abc9c', '#34495e'],
					label: 'Tag Stats'
				}],
				labels: []
			}

			$.getJSON(apiUrl.concat('dashboard/stats/tags')).then(function (response) {
				if (response.error) {alert(response.message); return}

				_.each(response.data, function (item, index) {
					canvas.data.datasets[0].data.push(item.count)
					canvas.data.labels.push(item.label)
				})

				if (canvas.chart) canvas.chart.destroy()
				canvas.chart = new Chart(canvas.context, {
					type: 'pie',
					data: canvas.data,
					options: {
						animation: {
							animateRotate: true
						}
					}
				})
			})
		},

		getBestSelling: function (date) {
			var self = this
			self.loader.bestselling = true
			self.tab.bestselling = date

			var endpoint = apiUrl.concat('dashboard/stats/market?bestselling=1&date=', date)
			$.getJSON(endpoint).then(function (response) {
				if (response.error) { alert(response.message); return }

				// Request image preview
				var data = []
				var queue = async.queue(function (task, callback) {
					task.title = task.title.replace("- WordPress | ThemeForest", "").replace("- WordPress", "")
					if (task.img_preview === null) {
						$.getJSON(apiUrl.concat('getPreview'), {uri: task.url}).then(function (response) {
							if (response.error) { alert(response.message); return}
							task.img_preview = response.img.url
							data.push(task)
							callback()
						})
					} else {
						data.push(task)
						callback()
					}
				});

				queue.drain = function () {
					data = _.chain(data).groupBy(function(element, index){
						return Math.floor(index/3);
					}).toArray().value()

					self.bestselling = data
					self.loader.bestselling = false
				};

				_.each(response.data, function (item) {
					queue.push(item)
				})

				if (_(data).size()===0) {
					self.bestselling = data
					self.loader.bestselling = false
				}
			})
		}
	},
	ready: function () {
		this.viewMarketValue('market-value', 'today')
		this.getCategories()
		this.getGroups()
		this.getTagStats()
		this.getBestSelling('today')

		this.$watch('selectedCategory', function () {
			this.viewMarketValue('categories-value', this.tab['categories-value'])
		})

		this.$watch('selectedGroup', function () {
			this.viewMarketValue('group-value', this.tab['group-value'])
		})
	}
});


/**
 * Application wrapper template
 * @type {Object}
 */
App.wrapper = Vue.extend({
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
		searchForm: {
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

App.groupSelector = Vue.extend({
	template: '#group-selector',
	props: {
		show: {
			required: true,
			type: Boolean,
			default: false
		},
		done: {
			type: Function,
			default: null
		}
	},
	data: function () {
		return {
			selectedList: '',
			list: []
		}
	},
	ready: function () {
		var self = this
		$.getJSON(apiUrl.concat('list/subscribe/group')).then(function (response) {
			if (response.error) {
				alert(response.message)
				return
			}

			if (response.list && response.list.length > 0) {
				self.selectedList = response.list[0]
				self.list = response.list
			}
		})
	},
	methods: {
		ok: function () {
			this.done && this.done(this.selectedList)
			this.show = false
		}
	}
});


/**
 * Application page components
 * @type {Object}
 */
App.page = {
	/** 
	* View all page list
	*/
	list: Vue.extend({
		template: '#watcher-list',
		components: {wrapper: App.wrapper},
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
		components: {wrapper: App.wrapper},
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
					alert("URL and Title seems not valid?")
				}
			}
		}
	}),


	/**
	 * Edit page
	 * @type {Object}
	 */
	edit: Vue.extend({
		components: {wrapper: App.wrapper},
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
};


/**
 * Application item component
 * @type {Object}
 */
App.item = {
	/**
	 * Item List
	 * @type {Object}
	 */
	list: Vue.extend({
		template: '#item-list',
		components: {wrapper: App.wrapper, 'group-selector': App.groupSelector},
		route: {
			canReuse: false,
			data: function () {
				return this.fetch()
			},
		},
		
		data: function () {
			return {
				bulkAction: '',
				allChecked: false,
				checkAllItem: false,
				checkedList: [],
				groupSelector: false,
				paginationLoaded: false,
				pagination: {items: 0, itemsOnPage: 100}
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

					self.pagination.items = response.total
					response.list = _.map(response.list, function (item, index) {
						item.created = renderDate(item.created)
						return item
					})
					return response
				})
			},


			changedPage: function (response) {
				this.list = response.list
				this.allChecked = false
				this.checkedList = []
				this.bulkAction = ''
			},

			getPage: function (offset) {
				this.fetch(offset).then(this.changedPage)
			},

			prevPage: function (event) {
				var offset = this.pagination.current - 1

				offset = (offset <= 0)? 0: offset
				this.fetch(offset).then(this.changedPage)
			},

			nextPage: function (event) {
				var offset = this.pagination.current + 1

				offset = (offset >= this.pagination.total)? this.pagination.total - 1: offset
				this.fetch(offset).then(this.changedPage)
			},

			checkAll: function (event) {
				this.checkedList = []
				this.bulkAction = ''

				if (!this.allChecked) {

					var confirmAllItems = confirm("There is " + this.total + " items, select all of them?")
					this.checkAllItem = confirmAllItems

					for (item in this.list) {
						this.checkedList.push(this.list[item].item_id);
					}
				}
			},

			applyBulkAction: function () {
				if (this.checkedList.length > 0) {
					if (this.bulkAction === 'subscribe') {
						this.groupSelector = true
					} else {
						this.subunsub(false)
					}
				}
			},

			selectedGroup: function (selected) {
				if (selected) this.subunsub(true, selected.id)
				else {
					this.bulkAction = ''
					this.checkedList = []
					this.allChecked = false
				}
			},

			subunsub: function (subscribe, group_id) {
				var self = this, type = (subscribe)? 'subscribe': 'unsubscribe'

				queues = async.queue(function (item_id, callback) {
					var params = {item_id: item_id}

					if (subscribe) params.group_id = group_id

					$.post(apiUrl.concat(type), params)
					.then(function (response) {

						if (!response.error) {
							_.findWhere(self.list, {item_id: item_id}).subscribed = subscribe
						}

						callback()
					})
				}, 5)

				queues.drain = function () {
					self.bulkAction = ''
					self.checkedList = []
					self.allChecked = false
				}

				_.each(_.unique(self.checkedList), function (item) {
					queues.push(item)
				})
			}
		},

		activate: function () {

		},

		ready: function () {
			var self = this

			if (self.$root && ! self.$root.paginationItemList) {
				var pagination = UIkit.pagination('#pagination', {
					displayedPages: 7,
					currentPage: 0
				})

				pagination.UIkit.on('select.uk.pagination', function(e, index) {
					self.getPage(index)
				})

				self.$root.paginationItemList = true
			}
		}
	}),


	/**
	 * Item View
	 * @type {Object}
	 */
	view: Vue.extend({
		template: '#item-view',
		components: {wrapper: App.wrapper, 'group-selector': App.groupSelector},
		route: {
			canReuse: false,
			waitForData: true,
			data: function () {
				return this.fetch()
			},
		},
		data: function () {
			return {
				tab: 'today',
				totalRevenue: 0,
				totalSales: 0,
				canvas: {},
				groupSelector: false,
				style: {}
			}
		},
		methods: {
			/**
			 * Get market value
			 * @param  {String} date
			 * @return {void}
			 */
			getMarketValue: function (date) {
				var self = this
				self.canvas.context = document.getElementById('item-chart').getContext('2d'),
				self.canvas.data = {
					labels: [],
					datasets: [{
						label: "Revenue",
						fill: true,
						lineTension: 0.1,
						backgroundColor: "rgba(26, 188, 156, 0.4)",
						borderColor: "rgba(26, 188, 156,1.0)",
						borderCapStyle: 'butt',
						borderDash: [],
						borderDashOffset: 0.0,
						borderJoinStyle: 'miter',
						pointBorderColor: "rgba(22, 160, 133,1.0)",
						pointBackgroundColor: "rgba(26, 188, 156,1.0)",
						pointBorderWidth: 2,
						pointHoverRadius: 5,
						pointHoverBackgroundColor: "#ffffff",
						pointHoverBorderColor: "rgba(22, 160, 133,1.0)",
						pointHoverBorderWidth: 2,
						pointRadius: 5,
						pointHitRadius: 10,
						data: [],
						sales: []
					}]
				}

				$.getJSON(apiUrl.concat('dashboard/stats/market?date=', date, '&item_id=', self.item_id))
				.then(function (response) {
					if (response.error) { alert(response.message); return }

					// Push actual data to chart
					_.each(orderDate(response.data), function (item, date) {
						self.canvas.data.labels.push(date)
						self.canvas.data.datasets[0].data.push(item.sales * item.price)
						self.canvas.data.datasets[0].sales.push(item.sales)
					})

					// Clear and Render Canvas
					if (self.canvas.chart) self.canvas.chart.destroy()
					self.canvas.chart = new Chart(self.canvas.context, {
						type: 'line',
						data: self.canvas.data,
						options: {
							 tooltips: {
							 	callbacks: {
							 		label: function(item, data) {
							 			var sales = data.datasets[item.datasetIndex].sales[item.index]
							 			return sales +' x ' + self.price + ' = $' + (item.yLabel).toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,')
							 		}
							 	}
							 }
						}
					})
				})
			},


			viewMarketValue: function (date) {
				this.getMarketValue(date)
				this.tab = date
			},

			getImgPreview: function () {
				var self = this
				$.getJSON(apiUrl.concat('getPreview'), {uri: self.url})
				.then(function (response) {
					if (response.error) {
						alert(response.message);
						return
					}

					self.$set('style', {
						backgroundImage: 'url('+ response.img.url +')'
					})
				})
			},

			fetch: function () {
				var self = this
				return $.getJSON(apiUrl.concat('item/', this.$route.params.id)).then(function (response) {
					if (response.error) {
						alert(response.message)
						return
					}

					self.totalSales = _.last(response.data.sales).value
					self.totalRevenue = (self.totalSales * response.data.price).toFixed(2).replace(/(\d)(?=(\d{3})+\.)/g, '$1,')
					return response.data
				})
			},

			subscribe: function () {
				if (!this.subscribed) this.groupSelector = true
				else this.subunsub(false)
			},

			subunsub: function (selected) {
				var self = this, type = (selected)? 'subscribe': 'unsubscribe'
				$.post(apiUrl.concat(type), {item_id: self.item_id, group_id: selected.id}).then(function (response) {
					if (response.error) {
						alert(response.message)
					}

					if (type === 'subscribe') self.subscribed = true
					else self.subscribed = false
				})
			}
		},

		ready: function () {
			this.getMarketValue('today', this.item_id)
			this.getImgPreview()
		}
	})
};


/**
 * Application subscribe component
 * @type {Object}
 */
App.subscribe = {}
App.subscribe.wrapper = Vue.extend({
	template: '#subscribe',
	components: {wrapper: App.wrapper}
});
	

App.subscribe.list = Vue.extend({
	template: '#subscribe-list',
	components: {wrapper: App.wrapper, 'subscribe-wrapper': App.subscribe.wrapper},
	route: {
		canReuse: false,
		waitForData: true,
		data: function () {
			return this.fetch()
		},
	},
	methods: {
		fetch: function () {
			return $.getJSON('http://localhost:8080/api/v1/list/subscribe/group', function (response) {
				if (response.error) { alert(response.message); return }

				var groups = {}
				_.each(response.list, function (item, index) {
					groups[item.id] = item
					groups[item.id].items = []
				})

				$.getJSON('http://localhost:8080/api/v1/list/subscribe', function (response) {
					if (response.error) { alert(response.message); return }

					_.each(response.list, function (item, index) {
						item.created = renderDate(item.created)
						_.each(item.subscribe_group_id, function (_item, _index) {
							groups[_item].items.push(item)
						})
					})
				})

				return groups
			})
		}
	}
});

App.subscribe.groupList = Vue.extend({
	template: '#subscribe-group-list',
	components: {wrapper: App.wrapper, 'subscribe-wrapper': App.subscribe.wrapper},
	route: {
		canReuse: false,
		waitForData: true,
		data: function () {
			return this.fetch()
		},
	},
	methods: {
		fetch: function (fn) {
			return $.getJSON(apiUrl.concat('list/subscribe/group')).then(function (response) {
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
					url: apiUrl.concat('subscribe/group/', id),
					type: 'DELETE',
					dataType: 'JSON',
					success: function (response) {
						if (response.error) {
							alert(response.message)
							return
						}

						self.$router.go('/')
						setTimeout(function () {
							self.$router.go('/subscribe/group')
						}, 10)
					}
				});
			}
		}
	}
});

App.subscribe.addGroup = Vue.extend({
	template: '#subscribe-group-add',
	components: {wrapper: App.wrapper, 'subscribe-wrapper': App.subscribe.wrapper},
	methods: {
		valid: function () {
			return this.name !== '' && this.desc !== ''
		},
		save: function () {
			var self = this

			if (self.valid()) {
				$.post(apiUrl.concat('subscribe/group'), {
					name: self.name,
					desc: self.desc,
				}).then(function (response) {
					if (response.error) {
						alert(response.message)
					}

					self.$router.go('/subscribe/group')
				})
			} else {
				alert("Name and Desc not valid")
			}
		}
	}
});

App.subscribe.editGroup = Vue.extend({
	template: '#subscribe-group-edit',
	components: {wrapper: App.wrapper, 'subscribe-wrapper': App.subscribe.wrapper},
	data: function () {
		return {
			name: '',
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
			return $.getJSON(apiUrl.concat('subscribe/group/', this.$route.params.id)).then(function (response) {
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
				url: apiUrl.concat('subscribe/group/', self.$route.params.id),
				type: 'PUT',
				data: {title: self.title, desc: self.desc},
				dataType: 'JSON',
				success: function (response) {
					if (response.error) {
						alert(response.message)
						return
					}

					self.$router.go('/subscribe/group')
				}
			});
		}
	}
});


App.search = Vue.extend({
	template: '#search',
	components: {wrapper: App.wrapper},
	route: {
		canReuse: false,
		waitForData: true,
		data: function () {
			return this.fetchCategory()
		},
	},
	data: function () {
		return {
			currentTab: 'search',
			waiting: false,
			defaultFilter: {
				add: 'include',
				field: 'title',
				criteria: {number: '$gte', author: '', tags: '', category: '', 'title': ''},
				date: 'today',
				value: ''
			},

			filters: [{
				add: 'include',
				field: 'title',
				criteria: {number: '$gte', author: '', tags: '', category: '', 'title': ''},
				date: 'today',
				value: 'Wordpress'
			}],

			list: null,
			total: 0,
			params: null,
			pagination: {items: 0, itemsOnPage: 100} 
		}
	},
	methods: {
		fetchCategory: function () {
			return $.getJSON(apiUrl.concat('list/category')).then(function (response) {
				if (response.error) {
					alert(response.message)
					return
				}

				return {categories: response.list}
			})
		},

		fetch: function (offset) {
			if (offset === undefined) return
			var self = this
			$.post(apiUrl.concat('search', '?offset=', offset), {search: self.params}).then(function (response) {
				if (response.error) {
					alert(response.message)
					return
				}

				self.waiting = false
				self.list = response.list
				self.total = self.pagination.items = response.total

				// Activate pagination
				if (self.$root && ! self.$root.paginationSearch) {
					var pagination = UIkit.pagination('#pagination-search', {
						displayedPages: 7,
						currentPage: 0
					})

					pagination.UIkit.on('select.uk.pagination', function(e, index) {
						self.fetch(index)
					})

					self.$root.paginationSearch = true
				}
			}, "json")
		},

		remove: function (filter, index) {
			if (index === 0) return
			this.filters.$remove(filter)
		},

		add: function () {
			var filter = JSON.stringify(this.defaultFilter)
			this.filters.push(JSON.parse(filter))
		},

		search: function () {
			var self = this
			self.waiting = true
			self.currentTab = 'result'

			async.waterfall([
				// Get all include params
				function (next) {
					var searchFilters = {
						price: {$not: {}},
						author: {regex: [], options: "", not: false},
						tags: {regex: [], options: "", not: false},
						category: {regex: [], options: "", not: false},
						title: {regex: [], options: "", not: false},
						sales: {}
					} 

					// Clone filters
					var copyFilters = $.extend(true, {}, self.filters)

					// Parsing value based on it's criteria 
					_.each(copyFilters, function (item, index) {

						// Skip blank value
						if (item.value === '') return

						// Store item field
						var field = searchFilters[item.field],

						// Convert string to number for sales / price
						intValue = parseInt(item.value),
						isInteger = !isNaN(intValue)

						switch (item.field) {
							case 'price':
								if (isInteger) {
									if (item.add === 'include') {
										if (!field[item.criteria.number]) field[item.criteria.number] = []
										field[item.criteria.number].push(intValue)
									} else {
										if (!field.$not[item.criteria.number]) field.$not[item.criteria.number] = []
										field.$not[item.criteria.number].push(intValue)
									}
								}
							break;

							case 'sales':
							break;

							case 'tags':
							case 'author':
							case 'title':
								var regex = '',
								options = 'i',
								not = false,
								value = item.value

								// Strip |
								value = value.replace(/\|/g, "\\|").replace(/\-/g, "\\-")

								// Get item field
								if (item.field === 'tags') {
									value = value.split(',').map($.trim)
								} else {
									value = [value]
								}

								switch (item.criteria[item.field]) {
									case '$eq':
										regex = '(^'+ value.join('|') +'$)'
										options = ''
									break;

									case '$ne':
										regex = '(^'+ value.join('|') +'$)'
										not = true
										options = ''
									break;

									case '^':
										regex = '(^'+ value.join('|') +')'
										break;

									case '$':
										regex = '('+ value.join('|') +'$)'
										break;

									case '!c':
										regex = '(' + value.join('|') + ')'
										not = true
										break;

									default:
										regex = '(' + value.join('|') + ')'
										break;
								}

								field.regex.push(regex)
								field.options = options
								field.not = not
							break;

							case 'category':
								field.regex.push($.trim(item.criteria.category))
							break;
						}
					})

					next(null, searchFilters)
				},

				// Convert object to mongodb query
				// It's basicly JSON but with multiple keys, duh
				function (searchFilters, next) {
					var sort = function(num) {return num;}

					// Sort price and sales
					searchFilters.price.$gte = _.first(_.sortBy(searchFilters.price.$gte, sort))
					searchFilters.price.$lte = _.first(_.sortBy(searchFilters.price.$lte, sort).reverse())
					searchFilters.price.$not.$gte = _.first(_.sortBy(searchFilters.price.$not.$gte, sort))
					searchFilters.price.$not.$lte = _.first(_.sortBy(searchFilters.price.$not.$lte, sort).reverse())

					// Join all regex
					searchFilters.category.regex = searchFilters.category.regex.join('|')
					searchFilters.title.regex = searchFilters.title.regex.join('|')
					searchFilters.author.regex = searchFilters.author.regex.join('|')
					searchFilters.tags.regex = searchFilters.tags.regex.join('|')

					// If item is empty delete it
					var cleanSearchFilters = JSON.parse(JSON.stringify(searchFilters))
					_.each(cleanSearchFilters, function (item, key) {
						switch (key) {
							case 'price':
								if (_.isEmpty(item.$not)) delete cleanSearchFilters[key].$not
							break;

							default:
								if (_.isEmpty(item)) delete cleanSearchFilters[key]
								if (item.regex !== undefined) {
									if (_.isEmpty(item.regex)) delete cleanSearchFilters[key]
								}
							break;
						}
					})

					// If nothing to filters 
					
					if (!_.isEmpty(cleanSearchFilters)) {
						var params = _.pick(cleanSearchFilters, function (item) {
							return ! _.isEmpty(item)
						})
						params = JSON.stringify(params)
						next(null, params)
					}
				}

				// Done
			], function (err, params) {
				self.params = params
				self.fetch(0)
			})
		}
	},
	ready: function () {
		var self = this

		this.$watch('currentTab', function (value) {
			if (value === 'search') {
				this.waiting = false
				this.list = null
				this.total = 0
			}
		})
	}
})

// Main App, we need non instance
App.init = Vue.extend({
	data: function () {
		return {
			navigation: [
				{icon: 'dashboard', 'path': '/', label: 'Dashboard'},
				{icon: 'rss', 'path': '/page', label: 'Watcher'},
				{icon: 'file-o', 'path': '/item', label: 'Item List'},
				{icon: 'bookmark-o', 'path': '/subscribe', label: 'Subscribes'},
				{icon: 'search', 'path': '/search', label: 'Advanced Search'},
			]
		}
	}
})

// Router
App.router = new VueRouter({
	linkActiveClass: 'active'
})


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
	},

	'/subscribe': {
		component: App.subscribe.list
	},

	'/subscribe/group': {
		component: App.subscribe.groupList
	},

	'/subscribe/group/add': {
		component: App.subscribe.addGroup
	},

	'/subscribe/group/:id': {
		component: App.subscribe.editGroup
	},

	'/search': {
		component: App.search
	}
});

// Start
App.router.start(App.init, '#app')