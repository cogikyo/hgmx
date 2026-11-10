import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector("rect");
  animate(screen, { scaleY: [0.1, 1.1, 1], opacity: [0.5, 1] }, { duration: τmid });

  const play = el.querySelector('path[d^="M15.033"]');
  animate(play, { translateX: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
