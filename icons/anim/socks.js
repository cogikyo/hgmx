import { animate, τmax } from "../motion.js";

export default function (el) {
  const socks = el.querySelector("g");
  animate(socks, { translateY: [0, -3, 1, 0], rotate: [0, 10, -4, 0] }, { duration: τmax });
}
