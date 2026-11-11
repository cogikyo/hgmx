import { animate, τmax } from "../motion.js";

export default function (el) {
  const glasses = el.querySelectorAll(":scope > *");
  animate(glasses, { translateX: [0, 2.5, -1, 0], translateY: [0, -2.5, 1, 0] }, { duration: τmax });
}
