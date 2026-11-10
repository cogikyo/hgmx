import { animate, τmax } from "../motion.js";

export default function (el) {
  const shirt = el.querySelector("path");
  animate(shirt, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
