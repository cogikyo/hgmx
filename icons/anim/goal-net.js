import { animate, τmax } from "../motion.js";

export default function (el) {
  const net = el.querySelector("g");
  animate(net, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });
}
