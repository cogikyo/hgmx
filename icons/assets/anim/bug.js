import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bug = el.querySelector("g");
  animate(bug, { rotate: [0, -12, 10, -6, 0] }, { duration: τmax });

  const antennae = el.querySelectorAll('path[d^="M14.12"], path[d^="m8 2"]');
  animate(antennae, { rotate: [0, -30, 20, 0] }, { duration: τmid, delay: τmin });
}
