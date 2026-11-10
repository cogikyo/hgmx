import { animate, τmax } from "../motion.js";

export default function (el) {
  const glasses = el.querySelectorAll(":scope > *");
  animate(glasses, { translateY: [0, -2, 1, 0] }, { duration: τmax });
}
