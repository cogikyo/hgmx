import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wallet = el.querySelector("rect");
  animate(wallet, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pocket = el.querySelector('path[d^="M3 11"]');
  animate(pocket, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
