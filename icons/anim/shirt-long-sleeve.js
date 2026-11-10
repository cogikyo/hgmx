import { animate, τmax } from "../motion.js";

export default function (el) {
  const shirt = el.querySelectorAll("path");
  animate(shirt, { translateY: [0, -2, 0.5, 0] }, { duration: τmax });
}
