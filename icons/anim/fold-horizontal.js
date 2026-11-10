import { animate, τmax } from "../motion.js";

export default function (el) {
  const left = el.querySelector("g.left");
  animate(left, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const right = el.querySelector("g.right");
  animate(right, { translateX: [4, -1, 0], opacity: [0, 1] }, { duration: τmax });
}
