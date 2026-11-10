import { animate, τmax } from "../motion.js";

export default function (el) {
  const blade = el.querySelector("g");
  animate(blade, { translateX: [0, 3, -1, 0] }, { duration: τmax });
}
