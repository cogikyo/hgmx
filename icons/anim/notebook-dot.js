import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cover = el.querySelectorAll('path[d^="M20 11.75"], path[d^="M16 11.75"]');
  animate(cover, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
