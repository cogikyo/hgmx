import { animate, τmax } from "../motion.js";

export default function (el) {
  const skate = el.querySelectorAll("path");
  animate(skate, { translateX: [0, 4, -1, 0] }, { duration: τmax });
}
