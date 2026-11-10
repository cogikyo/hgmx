import { animate, τmax } from "../motion.js";

export default function (el) {
  const whisk = el.querySelector("g");
  animate(whisk, { rotate: [0, -12, 4, 0] }, { duration: τmax });

  const knife = el.querySelector('path[d^="M22 13"]');
  animate(knife, { rotate: [0, 12, -4, 0] }, { duration: τmax });
}
