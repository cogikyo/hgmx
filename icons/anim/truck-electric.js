import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("g");
  const truck = [body, ...el.querySelectorAll("circle")];
  animate(truck, { translateX: [-6, 2, 0], opacity: [0, 1] }, { duration: τmax });

  const bolt = el.querySelector('path[d^="M4 3"]');
  animate(bolt, { scale: [1, 1.35, 1] }, { duration: τmid, delay: τmin });

  animate(body, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
