import { animate, τmax } from "../motion.js";

export default function (el) {
  const hat = el.querySelectorAll("path, rect");
  animate(hat, { translateY: [0, -2, 0.5, 0] }, { duration: τmax });
}
