import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const orbits = [el.querySelector('path[d^="M20.2"]'), el.querySelector('path[d^="M15.7"]')];
  animate(orbits, { rotate: [0, 180] }, { duration: τmax });

  const nucleus = el.querySelector("circle");
  animate(nucleus, { scale: [0, 1.75, 1] }, { duration: τmid, delay: τmin });
}
