import { animate, τmax } from "../motion.js";

export default function (el) {
  const pine = el.querySelector("g");
  animate(pine, { rotate: [0, -14, 8, 0] }, { duration: τmax });
}
