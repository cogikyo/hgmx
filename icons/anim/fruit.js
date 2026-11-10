import { animate, τmax } from "../motion.js";

export default function (el) {
  const fruit = el.querySelectorAll("circle, path");
  animate(fruit, { translateY: [-4, 1, 0] }, { duration: τmax });
}
