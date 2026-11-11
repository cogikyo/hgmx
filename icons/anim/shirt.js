import { animate, τmax } from "../motion.js";

export default function (el) {
  const shirt = el.querySelector("path");
  animate(shirt, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
