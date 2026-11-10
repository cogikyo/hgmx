import { animate, τmax } from "../motion.js";

export default function (el) {
  const hat = el.querySelector("g");
  animate(hat, { translateY: [0, -3, 0], rotate: [0, -12, 5, 0] }, { duration: τmax });
}
