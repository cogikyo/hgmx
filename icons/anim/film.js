import { animate, τmax } from "../motion.js";

export default function (el) {
  const strip = el.querySelector("g");
  animate(strip, { translateY: [-4.5, 1, 0] }, { duration: τmax });
}
