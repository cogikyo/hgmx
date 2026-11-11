import { animate, τmax } from "../motion.js";

export default function (el) {
  const shade = el.querySelector("path[transform-origin]");
  animate(shade, { rotate: [0, -55, 30, -12, 0] }, { duration: τmax });
}
