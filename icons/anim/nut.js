import { animate, τmax } from "../motion.js";

export default function (el) {
  const nut = el.querySelector("g");
  animate(nut, { translateY: [-4, 1, 0] }, { duration: τmax });
}
