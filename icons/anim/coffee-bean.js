import { animate, τmax } from "../motion.js";

export default function (el) {
  const bean = el.querySelector("g");
  animate(bean, { rotate: [0, -20, 10, 0], scale: [1, 1.15, 1] }, { duration: τmax });
}
