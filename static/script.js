const { createApp } = Vue;

const contact = createApp({
  methods: {
    submitForm: function (x) {
      var button = document.getElementById("button");
      button.innerText = "Sending...";
      button.disabled = true;
      fetch("/", {
        method: "POST",
        headers: {
          "X-CSRF-Token": x.target.elements["gorilla.csrf.Token"].value,
        },
        credentials: "same-origin",
        body: new FormData(x.target),
      }).then((resp) => {
        if (resp.ok) {
          button.innerText = "Sent!";
        } else {
          button.innerText = "Error: " + resp.statusText;
        }
      });
    },
  },
});
contact.config.compilerOptions.delimiters = ["${", "}"];
contact.mount("#contact");

const count = createApp({
  data() {
    return { regions: {} };
  },
  methods: {
    async goGet() {
      var response = await fetch("/rank");
      this.regions = await response.json();
    },
  },
  created: function () {
    this.goGet();
  },
});
count.config.compilerOptions.delimiters = ["${", "}"];
count.mount("#count");

const pricing = createApp({
  data() {
    return { requests: 100, outbound: 0.01 };
  },
  computed: {
    requestcost: function () {
      return (this.requests * 0.2) / 1000000;
    },
    durationcost: function () {
      return this.requests * 0.000000208;
    },
    gwcost: function () {
      return this.requests * 0.00000425;
    },
    outboundcost: function () {
      return this.outbound * 0.12;
    },
    total: function () {
      return (
        this.requestcost +
        this.durationcost +
        this.gwcost +
        this.outboundcost
      ).toFixed(2);
    },
  },
});
pricing.config.compilerOptions.delimiters = ["${", "}"];
pricing.mount("#pricing");
