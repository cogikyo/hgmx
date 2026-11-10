import { animate, τmax } from "../motion.js";

export default function (el) {
  const chair = el.querySelectorAll(":scope > *");
  animate(chair, { translateX: [0, 2, -1, 0] }, { duration: τmax });

  const wheel = el.querySelector("g");
  animate(wheel, { rotate: [0, 23, -11.5, 0] }, { duration: τmax });
}
