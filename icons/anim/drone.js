import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const rotors = [...el.querySelectorAll("path[transform-origin]")];
  animate([rotors[0], rotors[3]], { rotate: [0, 45, -10, 0] }, { duration: τmax });
  animate([rotors[1], rotors[2]], { rotate: [0, -45, 10, 0] }, { duration: τmax });

  const drone = el.querySelectorAll("path, rect");
  animate(drone, { translateY: [0, -1.5, 0] }, { duration: τmid, delay: τmin });
}
