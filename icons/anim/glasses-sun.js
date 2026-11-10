import { animate, τmax } from "../motion.js";

export default function (el) {
  const shades = el.querySelectorAll(":scope > *");
  animate(shades, { translateY: [0, 2.5, -0.5, 0] }, { duration: τmax });
}
