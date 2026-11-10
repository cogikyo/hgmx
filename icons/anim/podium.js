import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const podium = el.querySelector("g");
  animate(podium, { scaleY: [0, 1.1, 1] }, { duration: τmid });
  const one = el.querySelector('path[d^="M12 6"]');
  animate(one, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
