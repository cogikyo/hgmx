import { animate, τmax } from "../motion.js";

export default function (el) {
  const ship = el.querySelector("g");
  animate(ship, { rotate: [0, -8, 6, -3, 0] }, { duration: τmax });

  const sea = el.querySelector('path[d^="M2 20"]');
  animate(sea, { translateX: [0, 2, -1, 0] }, { duration: τmax });
}
