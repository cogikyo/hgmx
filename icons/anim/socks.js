import { animate, τmax } from "../motion.js";

export default function (el) {
  const socks = el.querySelector("g");
  animate(socks, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
