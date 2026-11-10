import { animate, τmax } from "../motion.js";

export default function (el) {
  const kettle = el.querySelector("g");
  animate(kettle, { rotate: [0, 10, -3, 0] }, { duration: τmax });
}
